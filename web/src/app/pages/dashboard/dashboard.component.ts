import { OnInit, Component, OnDestroy } from '@angular/core';
import { InstanceService } from '../../instances/instance.service';
import { StatusService } from 'src/app/status/status.service';
import { GlobalsService } from 'src/app/globals/globals.service';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
  selector: 'ngx-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})

export class DashboardComponent implements OnInit, OnDestroy {

  public data: any[];
  public displayedColumns: string[] = [];
  public daprHealthiness: string;
  public daprVersion: string;
  private intervalHandler;

  constructor(
    private instanceService: InstanceService,
    private statusService: StatusService,
    public globals: GlobalsService,
    private snackbar: MatSnackBar,
  ) { }

  ngOnInit() {
    this.getInstances();
    this.getControlPlaneData();
    if (this.globals.kubernetesEnabled) {
      this.displayedColumns = ['name', 'labels', 'status', 'age', 'selector'];
    } else {
      this.displayedColumns = ['name', 'age', 'actions'];
    }

    this.intervalHandler = setInterval(() => {
      this.getInstances();
      this.getControlPlaneData();
    }, 3000);
  }

  ngOnDestroy() {
    clearInterval(this.intervalHandler);
  }

  checkEnvironment() {
    this.globals.getSupportedEnvironments();
  }

  getInstances() {
    this.instanceService.getInstances().subscribe((data: any[]) => {
      this.data = data;
    });
  }

  getControlPlaneData(): void {
    this.statusService.getControlPlaneStatus().subscribe((data: any[]) => {
      this.daprHealthiness = data.every((service) => {
        return service.Healthy == 'True'
      }) ? 'Healthy' : 'Unhealthy';
      data.forEach(service => {
        this.daprVersion = service.Version;
      });
    });
  }

  showSnackbar(message: string) {
    this.snackbar.open(message, '', {
      duration: 2000,
    });
  }

  delete(id: string) {
    this.instanceService.deleteInstance(id).subscribe(() => {
      this.showSnackbar('Deleted Dapr instance with ID ' + id);
    }, error => {
      this.showSnackbar('Failed to remove Dapr instance with ID ' + id);
    });
  }
}
