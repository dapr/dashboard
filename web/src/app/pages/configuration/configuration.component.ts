import { Component, OnInit } from '@angular/core';
import { ConfigurationService } from 'src/app/configuration/configuration.service';
import { DaprConfiguration } from 'src/app/types/types';

@Component({
  selector: 'app-configuration',
  templateUrl: './configuration.component.html',
  styleUrls: ['./configuration.component.scss']
})
export class ConfigurationComponent implements OnInit {

  public config: DaprConfiguration[];
  public displayedColumns: string[] = ['name', 'tracing-enabled', 'mtls-enabled', 'mtls-workload-ttl', 'mtls-clock-skew', 'age', 'created'];

  constructor(private configurationService: ConfigurationService) { }

  ngOnInit(): void {
    this.getConfiguration();
  }

  getConfiguration(): void {
    this.configurationService.getConfigurations().subscribe((data: DaprConfiguration[]) => {
      this.config = data;
    });
  }
}
