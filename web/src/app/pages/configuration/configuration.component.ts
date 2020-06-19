import { Component, OnInit } from '@angular/core';
import { ConfigurationService } from '../../configuration/configuration.service';

@Component({
  selector: 'app-configuration',
  templateUrl: './configuration.component.html',
  styleUrls: ['./configuration.component.scss']
})
export class ConfigurationComponent implements OnInit {

  public config: any[];

  constructor(private configService: ConfigurationService) {}

  ngOnInit(): void {
    this.getConfiguration();
  }

  public components: any[];
  public componentsStatus: any[];
  public displayedColumns: string[] = [];

  getConfiguration() {
    this.configService.getConfiguration().subscribe((data: any[]) => {
      this.config = data;
    });
  }
}
