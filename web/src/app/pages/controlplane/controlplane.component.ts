import { Component, OnInit } from '@angular/core';
import { StatusService } from '../../status/status.service';

@Component({
  selector: 'app-controlplane',
  templateUrl: './controlplane.component.html',
  styleUrls: ['./controlplane.component.scss']
})
export class ControlPlaneComponent implements OnInit {

  public data: any[];
  public displayedColumns: string[] = ['name', 'namespace', 'healthy', 'status', 'version', 'age', 'created'];

  constructor(private statusService: StatusService) { }

  ngOnInit(): void {
    this.getControlPlaneStatus();
  }

  getControlPlaneStatus() {
    this.statusService.getControlPlaneStatus().subscribe((data: any[]) => {
      this.data = data;
    });
  }
}
