import { Component, OnInit } from '@angular/core';
import { ConfigurationsService } from 'src/app/configurations/configurations.service';
import { DaprConfigurationStatus } from 'src/app/types/types';
import { ScopesService } from 'src/app/scopes/scopes.service';

@Component({
  selector: 'app-configuration',
  templateUrl: './configuration.component.html',
  styleUrls: ['./configuration.component.scss']
})
export class ConfigurationComponent implements OnInit {

  public config: DaprConfigurationStatus[];
  public displayedColumns: string[] = ['name', 'tracing-enabled', 'mtls-enabled', 'mtls-workload-ttl', 'mtls-clock-skew', 'age', 'created'];
  private intervalHandler;

  constructor(
    private configurationService: ConfigurationsService,
    private scopesService: ScopesService
  ) { }

  ngOnInit(): void {
    this.getConfiguration();

    this.intervalHandler = setInterval(() => {
      this.getConfiguration();
    }, 3000);

    this.scopesService.scopeChanged.subscribe(() => {
      this.getConfiguration();
    });
  }

  getConfiguration(): void {
    this.configurationService.getConfigurationsStatus().subscribe((data: DaprConfigurationStatus[]) => {
      this.config = data;
    });
  }
}
