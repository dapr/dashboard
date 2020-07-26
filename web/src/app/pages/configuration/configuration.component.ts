import { Component, OnInit } from '@angular/core';
import { ConfigurationsService } from 'src/app/configurations/configurations.service';
import { DaprConfigurationStatus } from 'src/app/types/types';

@Component({
  selector: 'app-configuration',
  templateUrl: './configuration.component.html',
  styleUrls: ['./configuration.component.scss']
})
export class ConfigurationComponent implements OnInit {

  public config: DaprConfigurationStatus[];
  public displayedColumns: string[] = ['name', 'tracing-enabled', 'mtls-enabled', 'mtls-workload-ttl', 'mtls-clock-skew', 'age', 'created'];

  constructor(private configurationService: ConfigurationsService) { }

  ngOnInit(): void {
    this.getConfiguration();
  }

  getConfiguration(): void {
    this.configurationService.getConfigurationsStatus().subscribe((data: DaprConfigurationStatus[]) => {
      this.config = data;
    });
  }
}
