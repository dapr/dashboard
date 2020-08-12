import { Component, OnInit, OnDestroy } from '@angular/core';
import { ComponentsService } from 'src/app/components/components.service';
import { ScopesService } from 'src/app/scopes/scopes.service';
import { DaprComponent } from 'src/app/types/types';

@Component({
  selector: 'app-components',
  templateUrl: 'dapr-components.component.html',
  styleUrls: ['dapr-components.component.scss'],
})
export class DaprComponentsComponent implements OnInit, OnDestroy {

  public components: DaprComponent[];
  public componentsLoaded: boolean;
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
    }, 10000);

    this.scopesService.scopeChanged.subscribe(() => {
      this.getComponents();
    });
  }

  ngOnDestroy(): void {
    clearInterval(this.intervalHandler);
  }

  getComponents(): void {
    this.componentsService.getComponents().subscribe((data: DaprComponent[]) => {
      this.components = data;
      this.components.forEach(component => {
        component.img = this.getIconPath(component.type);
      });
      this.componentsLoaded = true;
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
    } else {
      return 'assets/images/secretstores.png';
    }
  }
}
