import { Component, OnInit } from '@angular/core';
import { InstanceService } from 'src/app/instances/instance.service';
import { Status } from 'src/app/types/types';
import { ScopesService } from 'src/app/scopes/scopes.service';

@Component({
  selector: 'app-controlplane',
  templateUrl: './controlplane.component.html',
  styleUrls: ['./controlplane.component.scss']
})
export class ControlPlaneComponent implements OnInit {

  public data: Status[];
  public displayedColumns: string[] = ['name', 'namespace', 'healthy', 'status', 'version', 'age', 'created'];
  private intervalHandler;

  constructor(
    private statusService: InstanceService,
    private scopesService: ScopesService
  ) { }

  ngOnInit(): void {
    this.getControlPlaneStatus();

    this.intervalHandler = setInterval(() => {
      this.getControlPlaneStatus();
    }, 3000);

    this.scopesService.scopeChanged.subscribe(() => {
      this.getControlPlaneStatus();
    });
  }

  getControlPlaneStatus(): void {
    this.statusService.getControlPlaneStatus().subscribe((data: Status[]) => {
      this.data = data.sort((a, b) => a.name.localeCompare(b.name));
    });
  }
}
