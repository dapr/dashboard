import { Component, OnInit, OnDestroy } from '@angular/core';
import { InstanceService } from 'src/app/instances/instance.service';
import { ActivatedRoute } from '@angular/router';
import * as yaml from 'js-yaml';
import { GlobalsService } from 'src/app/globals/globals.service';
import { Metadata, Instance } from 'src/app/types/types';
import { ThemeService } from 'src/app/theme/theme.service';
import { YamlViewerOptions } from 'src/app/types/types';

@Component({
  selector: 'app-detail',
  templateUrl: 'detail.component.html',
  styleUrls: ['detail.component.scss'],
})
export class DetailComponent implements OnInit, OnDestroy {

  private id: string;
  public model: string;
  public modelYAML: any;
  public annotations: string[];
  public instance: Instance;
  public loadedConfiguration: boolean;
  public loadedInstance: boolean;
  public loadedMetadata: boolean;
  public metadata: Metadata[];
  public metadataDisplayedColumns: string[] = ['type', 'count'];
  public options: YamlViewerOptions;
  public platform: string;
  private intervalHandler;

  constructor(
    private route: ActivatedRoute,
    private instanceService: InstanceService,
    public globals: GlobalsService,
    private themeService: ThemeService,
  ) { }

  ngOnInit(): void {
    this.id = this.route.snapshot.params.id;
    this.checkPlatform();
    this.loadData();
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

    this.intervalHandler = setInterval(() => {
      this.getMetadata(this.id);
    }, 10000);
  }

  ngOnDestroy(): void {
    clearInterval(this.intervalHandler);
  }

  getConfiguration(id: string): void {
    this.instanceService.getDeploymentConfiguration(id).subscribe((data: string) => {
      this.model = data;
      try {
        this.modelYAML = yaml.safeLoad(data);
        this.annotations = Object.keys(this.modelYAML.metadata.annotations);
        this.loadedConfiguration = true;
      } catch (e) {
        this.modelYAML = {};
      }
    });
  }

  checkPlatform(): void {
    this.globals.getPlatform().subscribe(platform => { this.platform = platform; });
  }

  getInstance(id: string): void {
    this.instanceService.getInstance(id).subscribe((data: Instance) => {
      this.instance = data;
      this.loadedInstance = true;
    });
  }

  getMetadata(id: string): void {
    this.instanceService.getMetadata(id).subscribe((data: Metadata[]) => {
      this.metadata = data;
      this.loadedMetadata = true;
    });
  }

  loadData(): void {
    this.getConfiguration(this.id);
    this.getInstance(this.id);
    this.getMetadata(this.id);
  }
}
