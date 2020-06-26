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
  public controlPlaneHealthiness: number;
  private intervalHandler;

  constructor(
    private instanceService: InstanceService,
    private statusService: StatusService
  ) { }

  ngOnInit() {
    this.getInstances();
    this.getControlPlaneHealthiness();
    
    this.intervalHandler = setInterval(() => { 
      this.getInstances();
      this.getControlPlaneHealthiness();
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

  getControlPlaneHealthiness(): void {
    var totalHealthy = 0;
    this.statusService.getControlPlaneStatus().subscribe((data: any[]) => {
      data.forEach(service => {
        if (service.Healthy == 'True') {
          totalHealthy += 1;
        }
      });
      this.controlPlaneHealthiness = totalHealthy;
    });
  }
}
