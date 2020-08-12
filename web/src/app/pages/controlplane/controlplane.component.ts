import { Component, OnInit, OnDestroy } from '@angular/core';
import { InstanceService } from 'src/app/instances/instance.service';
import { Status } from 'src/app/types/types';

@Component({
  selector: 'app-controlplane',
  templateUrl: './controlplane.component.html',
  styleUrls: ['./controlplane.component.scss']
})
export class ControlPlaneComponent implements OnInit, OnDestroy {

  public data: Status[];
  public displayedColumns: string[] = ['name', 'namespace', 'healthy', 'status', 'version', 'age', 'created'];
  public controlPlaneLoaded: boolean;
  private intervalHandler;

  constructor(
    private statusService: InstanceService,
  ) { }

  ngOnInit(): void {
    this.getControlPlaneStatus();

    this.intervalHandler = setInterval(() => {
      this.getControlPlaneStatus();
    }, 10000);
  }

  ngOnDestroy(): void {
    clearInterval(this.intervalHandler);
  }

  getControlPlaneStatus(): void {
    this.controlPlaneLoaded = false;
    this.statusService.getControlPlaneStatus().subscribe((data: Status[]) => {
      this.data = data.sort((a, b) => a.name.localeCompare(b.name));
      this.controlPlaneLoaded = true;
    });
  }
}
