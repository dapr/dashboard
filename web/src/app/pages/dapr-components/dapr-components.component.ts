import { Component, OnInit } from '@angular/core';
import { ComponentsService } from 'src/app/components/components.service';
import { DaprComponentStatus } from 'src/app/types/types';
import { ScopesService } from 'src/app/scopes/scopes.service';

@Component({
  selector: 'app-components',
  templateUrl: 'dapr-components.component.html',
  styleUrls: ['dapr-components.component.scss'],
})
export class DaprComponentsComponent implements OnInit {

  public componentsStatus: DaprComponentStatus[];
  public statusLoaded: boolean;
  public displayedColumns: string[] = ['img', 'name', 'status', 'age', 'created'];
  private intervalHandler;

  constructor(
    private componentsService: ComponentsService,
    private scopesService: ScopesService
  ) { }

  ngOnInit(): void {
    this.getComponents();

    this.intervalHandler = setInterval(() => {
      this.getComponents();
    }, 3000);

    this.scopesService.scopeChanged.subscribe(() => {
      this.getComponents();
    });
  }

  getComponents(): void {
    this.statusLoaded = false;
    this.componentsService.getComponentsStatus().subscribe((data: DaprComponentStatus[]) => {
      this.componentsStatus = data;
      this.componentsStatus.forEach(component => {
        component.img = this.getIconPath(component.type);
      });
      this.statusLoaded = true;
    });
  }

  getIconPath(type: string): string {
    if (type.includes('bindings')) {
      return 'assets/images/bindings.png';
    } else if (type.includes('secretstores')) {
      return 'assets/images/secretstores.png';
    } else if (type.includes('state')) {
      return 'assets/images/statestores.png';
    } else if (type.includes('pubsub')) {
      return 'assets/images/pubsub.png';
    } else if (type.includes('exporters')) {
      return 'assets/images/tracing.png';
    }
  }
}
