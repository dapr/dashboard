import { Component, OnInit, OnDestroy } from '@angular/core';
import { ConfigurationsService } from 'src/app/configurations/configurations.service';
import { ScopesService } from 'src/app/scopes/scopes.service';
import { DaprConfiguration } from 'src/app/types/types';
import { GlobalsService } from 'src/app/globals/globals.service';

@Component({
  selector: 'app-configuration',
  templateUrl: './configuration.component.html',
  styleUrls: ['./configuration.component.scss']
})
export class ConfigurationComponent implements OnInit, OnDestroy {

  public config: DaprConfiguration[];
  public displayedColumns: string[] = [];
  public configurationsLoaded: boolean;
  private intervalHandler;

  constructor(
    private configurationService: ConfigurationsService,
    public globals: GlobalsService,
    private scopesService: ScopesService,
  ) { }

  ngOnInit(): void {
    this.globals.getPlatform().subscribe(platform => {
      if (platform === 'kubernetes') {
        this.displayedColumns = ['name', 'tracing-enabled', 'mtls-enabled', 'mtls-workload-ttl', 'mtls-clock-skew', 'age', 'created'];
      } else if (platform === 'standalone') {
        this.displayedColumns = ['name', 'tracing-enabled', 'mtls-enabled', 'age', 'created'];
      }
    });

    this.getConfiguration();

    this.intervalHandler = setInterval(() => {
      this.getConfiguration();
    }, 10000);

    this.scopesService.scopeChanged.subscribe(() => {
      this.getConfiguration();
    });
  }

  ngOnDestroy(): void {
    clearInterval(this.intervalHandler);
  }

  getConfiguration(): void {
    this.configurationService.getConfigurations().subscribe((data: DaprConfiguration[]) => {
      this.config = data;
      this.configurationsLoaded = true;
    });
  }
}
