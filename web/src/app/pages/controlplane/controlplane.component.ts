import { Component, OnInit } from '@angular/core';
import { InstanceService } from 'src/app/instances/instance.service';
import { Status } from 'src/app/types/types';

@Component({
  selector: 'app-controlplane',
  templateUrl: './controlplane.component.html',
  styleUrls: ['./controlplane.component.scss']
})
export class ControlPlaneComponent implements OnInit {

  public data: Status[];
  public displayedColumns: string[] = ['name', 'namespace', 'healthy', 'status', 'version', 'age', 'created'];

  constructor(private statusService: InstanceService) { }

  ngOnInit(): void {
    this.getControlPlaneStatus();
  }

  getControlPlaneStatus(): void {
    this.statusService.getControlPlaneStatus().subscribe((data: Status[]) => {
      this.data = data;
    });
  }
}
