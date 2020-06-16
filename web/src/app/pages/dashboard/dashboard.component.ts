import { Component, OnDestroy } from '@angular/core';
import { InstanceService } from '../../instances/instance.service';

@Component({
  selector: 'ngx-dashboard',
  templateUrl: './dashboard.component.html'
})
export class DashboardComponent implements OnDestroy {
  public data: any[];
  displayedColumns: string[] = ['name'];
  private intervalHandler;

  constructor(private instanceService: InstanceService) {
    this.getInstances();
    this.intervalHandler = setInterval(() => { this.getInstances() }, 3000);
  }

  getInstances() {
    this.instanceService.getInstances().subscribe((data: any[]) => {
      this.data = data;
    });
  }

  ngOnDestroy() {
    clearInterval(this.intervalHandler);
  }
}
