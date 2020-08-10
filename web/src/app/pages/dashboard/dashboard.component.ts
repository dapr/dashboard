import { OnInit, Component, OnDestroy } from '@angular/core';
import { InstanceService } from 'src/app/instances/instance.service';
import { GlobalsService } from 'src/app/globals/globals.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Instance, Status } from 'src/app/types/types';

@Component({
  selector: 'app-dashboard',
  templateUrl: 'dashboard.component.html',
  styleUrls: ['dashboard.component.scss'],
})

export class DashboardComponent implements OnInit, OnDestroy {

  public instances: Instance[];
  public displayedColumns: string[] = [];
  public daprHealthiness: string;
  public daprVersion: string;
  public tableLoaded: boolean;
  public controlPlaneLoaded: boolean;
  private intervalHandler;

  constructor(
    private instanceService: InstanceService,
    public globals: GlobalsService,
    private snackbar: MatSnackBar,
  ) { }

  ngOnInit(): void {
    this.getInstances();
    this.getControlPlaneData();
    this.checkPlatform();

    this.intervalHandler = setInterval(() => {
      this.getInstances();
      this.getControlPlaneData();
    }, 3000);
  }

  ngOnDestroy(): void {
    clearInterval(this.intervalHandler);
  }

  checkPlatform(): void {
    this.globals.getPlatform().subscribe(platform => {
      if (platform === 'kubernetes') {
        this.globals.kubernetesEnabled = true;
        this.displayedColumns = ['name', 'labels', 'status', 'age', 'selector'];
      }
      else if (platform === 'standalone') {
        this.globals.standaloneEnabled = true;
        this.displayedColumns = ['name', 'age', 'actions'];
      }
      this.tableLoaded = true;
    });
  }

  getInstances(): void {
    this.instanceService.getInstances().subscribe((data: Instance[]) => {
      this.instances = data;
    });
  }

  getControlPlaneData(): void {
    this.instanceService.getControlPlaneStatus().subscribe((data: Status[]) => {
      this.daprHealthiness = data.every((service) => {
        return service.healthy === 'True';
      }) ? 'Healthy' : 'Unhealthy';
      data.forEach(service => {
        this.daprVersion = service.version;
      });
      this.controlPlaneLoaded = true;
    });
  }

  showSnackbar(message: string): void {
    this.snackbar.open(message, '', {
      duration: 2000,
    });
  }

  delete(id: string): void {
    this.instanceService.deleteInstance(id).subscribe(() => {
      this.showSnackbar('Deleted Dapr instance with ID ' + id);
    }, error => {
      this.showSnackbar('Failed to remove Dapr instance with ID ' + id);
    });
  }
}
