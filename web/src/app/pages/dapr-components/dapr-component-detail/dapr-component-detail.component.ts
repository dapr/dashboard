import { Component, OnInit } from '@angular/core';
import { DaprComponent } from 'src/app/types/types';
import { ComponentsService } from 'src/app/components/components.service';
import { ActivatedRoute } from '@angular/router';
import * as yaml from 'js-yaml';
import { ThemeService } from 'src/app/theme/theme.service';

@Component({
  selector: 'app-dapr-component-detail',
  templateUrl: './dapr-component-detail.component.html',
  styleUrls: ['./dapr-component-detail.component.scss']
})
export class DaprComponentDetailComponent implements OnInit {

  private name: string | undefined;
  public component: any;
  public componentManifest!: string;
  public loadedComponent = false;

  constructor(
    private route: ActivatedRoute,
    private componentsService: ComponentsService,
    private themeService: ThemeService,
  ) { }

  ngOnInit(): void {
    this.name = this.route.snapshot.params.name;
    if (typeof this.name !== 'undefined') {
      this.getComponent(this.name);
    }
  }

  getComponent(name: string): void {
    this.componentsService.getComponent(name).subscribe((data: DaprComponent) => {
      this.component = data;
      this.componentManifest = (typeof data.manifest === 'string') ?
        data.manifest : yaml.dump(data.manifest);
      this.loadedComponent = true;
    });
  }

  isDarkTheme() {
    return this.themeService.isDarkTheme();
  }
}
