import { Component, OnInit } from '@angular/core';
import { DaprComponent } from 'src/app/types/types';
import { ComponentsService } from 'src/app/components/components.service';
import { ActivatedRoute } from '@angular/router';
import * as yaml from 'js-yaml';
import { ThemeService } from 'src/app/theme/theme.service';
import { YamlViewerOptions } from 'src/app/types/types';

@Component({
  selector: 'app-dapr-component-detail',
  templateUrl: './dapr-component-detail.component.html',
  styleUrls: ['./dapr-component-detail.component.scss']
})
export class DaprComponentDetailComponent implements OnInit {

  private name: string;
  public component: any;
  public componentManifest: string | object;
  public options: YamlViewerOptions;
  public loadedComponent: boolean;

  constructor(
    private route: ActivatedRoute,
    private componentsService: ComponentsService,
    private themeService: ThemeService,
  ) { }

  ngOnInit(): void {
    this.name = this.route.snapshot.params.name;
    this.getComponent(this.name);
    this.options = {
      folding: true,
      minimap: { enabled: true },
      readOnly: false,
      language: 'yaml',
      theme: this.themeService.getTheme().includes('dark') ? 'vs-dark' : 'vs',
    };
    this.themeService.themeChanged.subscribe((newTheme: string) => {
      this.options = {
        ...this.options,
        theme: newTheme.includes('dark') ? 'vs-dark' : 'vs',
      };
    });
  }

  getComponent(name: string): void {
    this.componentsService.getComponent(name).subscribe((data: DaprComponent) => {
      this.component = data;
      this.componentManifest = yaml.safeDump(data.manifest);
      this.loadedComponent = true;
    });
  }
}
