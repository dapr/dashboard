import { Component, OnInit } from '@angular/core';
import { DaprConfiguration } from 'src/app/types/types';
import { ConfigurationsService } from 'src/app/configurations/configurations.service';
import { ActivatedRoute } from '@angular/router';
import * as yaml from 'js-yaml';
import { ThemeService } from 'src/app/theme/theme.service';
import { YamlViewerOptions, Instance } from 'src/app/types/types';
import { GlobalsService } from 'src/app/globals/globals.service';

@Component({
  selector: 'app-configuration-detail',
  templateUrl: './configuration-detail.component.html',
  styleUrls: ['./configuration-detail.component.scss']
})
export class ConfigurationDetailComponent implements OnInit {

  private name: string;
  public configuration: any;
  public configurationManifest: string | object;
  public options: YamlViewerOptions;
  public loadedConfiguration: boolean;
  public loadedApps: boolean;
  public configurationApps: Instance[];
  public platform: string;

  constructor(
    private route: ActivatedRoute,
    private configurationsService: ConfigurationsService,
    private themeService: ThemeService,
    public globals: GlobalsService,
  ) { }

  ngOnInit(): void {
    this.name = this.route.snapshot.params.name;
    this.checkPlatform();
    this.getConfiguration(this.name);
    this.getConfigurationApps(this.name);
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

  checkPlatform(): void {
    this.globals.getPlatform().subscribe(platform => { this.platform = platform; } );
  }

  getConfiguration(name: string): void {
    this.configurationsService.getConfiguration(name).subscribe((data: DaprConfiguration) => {
      this.configuration = data;
      this.configurationManifest = yaml.safeDump(data.manifest);
      this.loadedConfiguration = true;
    });
  }

  getConfigurationApps(name: string): void {
    this.configurationsService.getConfigurationApps(name).subscribe((data: Instance[]) => {
      this.configurationApps = data;
      this.loadedApps = true;
    });
  }
}
