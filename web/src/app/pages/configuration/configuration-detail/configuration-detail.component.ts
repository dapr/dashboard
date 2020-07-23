import { Component, OnInit } from '@angular/core';
import { DaprConfiguration } from 'src/app/types/types';
import { ConfigurationsService } from 'src/app/configurations/configurations.service';
import { ActivatedRoute } from '@angular/router';
import * as yaml from 'js-yaml';

@Component({
  selector: 'app-configuration-detail',
  templateUrl: './configuration-detail.component.html',
  styleUrls: ['./configuration-detail.component.scss']
})
export class ConfigurationDetailComponent implements OnInit {

  private name: string;
  public configuration: any;
  public configurationMetadata: string | object;
  public configurationDeployment: string | object;
  public options: object;
  public loadedConfiguration: boolean;

  constructor(
    private route: ActivatedRoute,
    private configurationsService: ConfigurationsService
  ) { }

  ngOnInit(): void {
    this.name = this.route.snapshot.params.name;
    this.getConfiguration(this.name);
    this.options = {
      folding: true,
      minimap: { enabled: true },
      readOnly: false,
      language: 'yaml',
    };
  }

  getConfiguration(name: string): void {
    this.configurationsService.getConfiguration(name).subscribe((data: DaprConfiguration) => {
      console.log(data);
      this.configuration = data;
      this.configurationDeployment = yaml.safeDump(data);
      this.configurationMetadata = yaml.safeDump(data.spec);
      this.loadedConfiguration = true;
    });
  }
}
