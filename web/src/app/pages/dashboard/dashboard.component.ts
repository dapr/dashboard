import { OnInit, Component, OnDestroy } from '@angular/core';
import { InstanceService } from 'src/app/instances/instance.service';
import { GlobalsService } from 'src/app/globals/globals.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Instance, Status } from 'src/app/types/types';
import { ScopesService } from 'src/app/scopes/scopes.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: 'dashboard.component.html',
  styleUrls: ['dashboard.component.scss'],
})

export class DashboardComponent implements OnInit, OnDestroy {

  public instances: Instance[] = [];
  public displayedColumns: string[] = [];
  public daprHealthiness: 'Healthy' | 'Unhealthy' = 'Unhealthy';
  public daprVersion = '';
  public tableLoaded = false;
  public controlPlaneLoaded = false;
  public platform = '';
  private intervalHandler: any;

  constructor(
    private instanceService: InstanceService,
    public globals: GlobalsService,
    private snackbar: MatSnackBar,
    private scopesService: ScopesService,
  ) { }

  ngOnInit(): void {
    this.checkPlatform();
    this.loadData();

    this.intervalHandler = setInterval(() => {
      this.loadData();
    }, 10000);

    this.scopesService.scopeChanged.subscribe(() => {
      this.loadData();
    });
  }

  ngOnDestroy(): void {
    clearInterval(this.intervalHandler);
  }

  checkPlatform(): void {
    this.globals.getPlatform().subscribe(platform => {
      this.platform = platform;
      if (platform === 'kubernetes') {
        this.displayedColumns = ['name', 'labels', 'status', 'age', 'selector'];
      }
      else if (platform === 'standalone') {
        this.displayedColumns = ['name', 'age', 'actions'];
      }
      else if (platform === 'docker-compose') {
        this.displayedColumns = ['name', 'age'];
      }
    });
  }

  getInstances(): void {
    this.instanceService.getInstances().subscribe((data: Instance[]) => {
      this.instances = data;
      this.tableLoaded = true;
    });
  }

  getControlPlaneData(): void {
    this.instanceService.getControlPlaneStatus().subscribe((data: Status[]) => {
      this.daprHealthiness = data.every((service) => service.healthy === 'True') ? 'Healthy' : 'Unhealthy';
      if (data.length === 0) {
        this.daprHealthiness = 'Unhealthy';
      }
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

  stop(id: string): void {
    if (confirm('You are about to stop Dapr instance with ID "' + id + '". Confirm?')) {
      this.instanceService.stopInstance(id).subscribe(() => {
        this.showSnackbar('Stopped Dapr instance with ID ' + id);
        this.getInstances();
      }, error => {
        this.showSnackbar('Failed to stop Dapr instance with ID ' + id);
      });
    }
  }

  loadData(): void {
    this.getInstances();
    this.getControlPlaneData();
  }
}
