# Adding a new platform

Dashboard currently supports 3 platforms: Kubernetes, self-hosted and docker-compose (August 2022).


## Backend

Platform definitions are contained in [platforms.go](../../pkg/platforms/platforms.go). When adding a new platform, define a new constant for the platform in `platforms.go`.

If the new platform requires configuration arguments you can define optional go flags in [dashboard.go](../../cmd/dashboard.go) and pass them to `RunWebServer(...)` in [webserver.go](../../cmd/webserver.go).

The runtime application platform is defined in `cmd/webserver.go:RunWebServer()`. To add a new platform, logic needs to be defined here to determine how Dashboard will recognize it. Any API clients or other necessary structures should be passed as arguments to the constructor of `instances.NewInstances(...)` and the other backend struct constructors along with the platform.

In `pkg/instances.go`, `pkg/components.go`, and `pkg/configurations.go`, new methods should be defined for each new platform, following the current pattern. In these files, functions such as `GetInstance(scope string, id string)` and `GetScopes()` are defined. These abstracted functions will call the correct platform-specific function:

```go
// GetInstances returns the result of the appropriate environment's GetInstance function
func (i *instances) GetInstances(scope string) []Instance {
	return i.getInstancesFn(scope)
}
```

Where `i.getInstanceFn` is defined in the constructor as `getPlatformInstances`, e.g. `getKubernetesInstances` or `getStandaloneInstances` according to the platform supplied.

If the new platform supports a feature (e.g. components) make sure you update the `func (c *type) Supported() bool` method in the corresponding package to include the platform definition.  For example, in [components.go](../../pkg/components/components.go)

```go
// Supported checks whether or not the supplied platform is able to access Dapr components
func (c *components) Supported() bool {
	return c.platform == platforms.Kubernetes || c.platform == platforms.Standalone || c.platform == platforms.DockerCompose
}
```

## Frontend

For platform-specific features of Dashboard, there is a service defined in `web/src/app/globals/globals.service.ts` that retrieves platform information from the backend. Any component that needs to know the current platform can make a call to this service and handle the returned data appropriately:

```typescript
checkPlatform(): void {
  this.globals.getPlatform().subscribe(platform => {
    this.platform = platform;
    if (platform === 'kubernetes') {
      ...
    }
    else if (platform === 'standalone') {
      ...
    }
  });
}
```

For entire sections that should only be shown on certain platforms, the variable set from the above can be used in the template files. For example, in `dashboard.component.html`, the control plane should only be shown in kubernetes mode. So the entire element should be wrapped in a `ng-container` tag with a `*ngIf` checking the platform:

```html
<ng-container *ngIf="platform === 'kubernetes'">
  <h3 class="card-header">Dapr Control Plane</h3>
  <mat-card class="card-tiny mat-elevation-z8">
    ...
  </mat-card>
</ng-container>
```
