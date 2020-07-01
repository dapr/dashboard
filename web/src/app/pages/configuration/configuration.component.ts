import { Component, OnInit } from '@angular/core';
import { ConfigurationService } from '../../configuration/configuration.service';

@Component({
  selector: 'app-configuration',
  templateUrl: './configuration.component.html',
  styleUrls: ['./configuration.component.scss']
})
export class ConfigurationComponent implements OnInit {

  public config: any[];
  public displayedColumns: string[] = ['name', 'tracing-enabled', 'mtls-enabled', 'mtls-workload-ttl', 'mtls-clock-skew', 'age', 'created'];

  constructor(private configService: ConfigurationService) { }

  ngOnInit(): void {
    this.getConfiguration();
  }

  getConfiguration() {
    this.configService.getConfiguration().subscribe((data: any[]) => {
      this.config = data;
    });
  }
}
