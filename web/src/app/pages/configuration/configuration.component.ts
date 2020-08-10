import { Component, OnInit } from '@angular/core';
import { ConfigurationsService } from 'src/app/configurations/configurations.service';
import { DaprConfiguration } from 'src/app/types/types';
import { GlobalsService } from 'src/app/globals/globals.service';

@Component({
  selector: 'app-configuration',
  templateUrl: './configuration.component.html',
  styleUrls: ['./configuration.component.scss']
})
export class ConfigurationComponent implements OnInit {

  public config: DaprConfiguration[];
  public displayedColumns: string[];
  public configurationsLoaded: boolean;

  constructor(
    private configurationService: ConfigurationsService,
    public globalsService: GlobalsService,
  ) { }

  ngOnInit(): void {
    this.getConfiguration();
    if (this.globalsService.kubernetesEnabled) {
      this.displayedColumns = ['name', 'tracing-enabled', 'mtls-enabled', 'mtls-workload-ttl', 'mtls-clock-skew', 'age', 'created'];
    } else {
      this.displayedColumns = ['name', 'tracing-enabled', 'mtls-enabled', 'age', 'created'];
    }
  }

  getConfiguration(): void {
    this.configurationService.getConfigurations().subscribe((data: DaprConfiguration[]) => {
      this.config = data;
      this.configurationsLoaded = true;
    });
  }
}
