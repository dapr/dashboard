import { OnInit, Component, OnDestroy } from '@angular/core';
import { InstanceService } from '../../instances/instance.service';
import { StatusService } from 'src/app/status/status.service';

@Component({
  selector: 'ngx-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})

export class DashboardComponent implements OnInit, OnDestroy {

  public data: any[];
  public displayedColumns: string[] = ['name', 'status', 'age'];
  public daprHealthiness: string;
  public daprVersion: string;
  private intervalHandler;

  constructor(
    private instanceService: InstanceService,
    private statusService: StatusService
  ) { }

  ngOnInit() {
    this.getInstances();
    this.getControlPlaneData();

    this.intervalHandler = setInterval(() => {
      this.getInstances();
      this.getControlPlaneData();
    }, 3000);
  }

  ngOnDestroy() {
    clearInterval(this.intervalHandler);
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
}
