import { OnInit, Component, OnDestroy } from '@angular/core';
import { InstanceService } from '../../instances/instance.service';

@Component({
  selector: 'ngx-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})

export class DashboardComponent implements OnInit, OnDestroy {
  public data: any[];
  public displayedColumns: string[] = ['name', 'status', 'age'];
  private intervalHandler;

  constructor(private instanceService: InstanceService) { }

  getInstances() {
    this.instanceService.getInstances().subscribe((data: any[]) => {
      this.data = data;
    });
  }

  ngOnInit() {
    this.getInstances();
    this.intervalHandler = setInterval(() => { this.getInstances(); }, 3000);
  }

  ngOnDestroy() {
    clearInterval(this.intervalHandler);
  }
}
