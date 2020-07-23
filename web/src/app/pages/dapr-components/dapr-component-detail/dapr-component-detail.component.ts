import { Component, OnInit } from '@angular/core';
import { DaprComponent } from 'src/app/types/types';
import { ComponentsService } from 'src/app/components/components.service';
import { ActivatedRoute } from '@angular/router';
import * as yaml from 'js-yaml';

@Component({
  selector: 'app-dapr-component-detail',
  templateUrl: './dapr-component-detail.component.html',
  styleUrls: ['./dapr-component-detail.component.scss']
})
export class DaprComponentDetailComponent implements OnInit {

  private name: string;
  public component: any;
  public componentName: string;
  public componentMetadata: string | object;
  public componentDeployment: string | object;
  public options: object;
  public loadedComponent: boolean;

  constructor(
    private route: ActivatedRoute,
    private componentsService: ComponentsService
  ) { }

  ngOnInit(): void {
    this.name = this.route.snapshot.params.name;
    this.getComponent(this.name);
    this.options = {
      folding: true,
      minimap: { enabled: true },
      readOnly: false,
      language: 'yaml',
    };
  }

  getComponent(name: string): void {
    this.componentsService.getComponent(name).subscribe((data: DaprComponent) => {
      this.component = data;
      this.componentName = data.metadata.name;
      this.componentDeployment = yaml.safeDump(data);
      this.componentMetadata = yaml.safeDump(data.spec.metadata);
      this.loadedComponent = true;
    });
  }
}
