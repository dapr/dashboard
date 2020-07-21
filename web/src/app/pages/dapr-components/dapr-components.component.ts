import { Component, OnInit } from '@angular/core';
import { ComponentsService } from 'src/app/components/component.service';

@Component({
  selector: 'ngx-dashboard',
  templateUrl: 'dapr-components.component.html',
  styleUrls: ['dapr-components.component.scss'],
})
export class DaprComponentsComponent implements OnInit {
  
  public components: any[];
  public componentsStatus: any[];
  public statusLoaded: boolean;
  public displayedColumns: string[] = ['img', 'name', 'status', 'age', 'created'];

  constructor(private componentsService: ComponentsService) { }

  ngOnInit(): void {
    this.getComponents();
  }

  getComponents(): void {
    this.statusLoaded = false;
    this.componentsService.getComponents().subscribe((data: any[]) => {
      this.components = data;

      for (const c of this.components) {
        c.iconPath = this.getIconPath(c.spec.type);
        c.spec.metadata = JSON.stringify(c.spec.metadata, null, 2);
      }
    });
    this.componentsService.getComponentsStatus().subscribe((data: any[]) => {
      this.componentsStatus = data;
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
