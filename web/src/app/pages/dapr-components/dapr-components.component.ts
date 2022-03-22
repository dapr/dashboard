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

  public components: DaprComponent[] = [];
  public componentsLoaded = false;
  public displayedColumns: string[] = ['img', 'name', 'status', 'age', 'created'];
  private intervalHandler: any;

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
    // Each component has a type field with the following format:
    //  <component_type>.<variable_part>
    // So we can determine which type of component we are handling based 
    // on the first part of the "type" field. 
    const componentType = type.split(".").shift();

    switch (componentType) {
      case 'bindings':
        return 'cloud';
      case 'secretstores':
        return 'lock';
      case 'state':
        return 'storage';
      case 'pubsub':
        return 'send';
      case 'exporters':
        return 'account_tree';

      default:
        console.warn("Unknown component type:", componentType);
        return 'question_mark';
    }
  }
}
