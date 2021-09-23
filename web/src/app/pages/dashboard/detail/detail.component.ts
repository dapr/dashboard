import { Component, OnInit, OnDestroy } from '@angular/core';
import { InstanceService } from 'src/app/instances/instance.service';
import { ActivatedRoute } from '@angular/router';
import * as yaml from 'js-yaml';
import { GlobalsService } from 'src/app/globals/globals.service';
import { Metadata, Instance } from 'src/app/types/types';
import { ThemeService } from 'src/app/theme/theme.service';

@Component({
  selector: 'app-detail',
  templateUrl: 'detail.component.html',
  styleUrls: ['detail.component.scss'],
})
export class DetailComponent implements OnInit, OnDestroy {

  private id: string | undefined;
  public model!: string;
  public modelYAML: any;
  public annotations: string[] = [];
  public instance!: Instance;
  public loadedConfiguration = false;
  public loadedInstance = false;
  public loadedMetadata = false;
  public metadata!: Metadata;
  public platform!: string;
  private intervalHandler: any;

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

    this.intervalHandler = setInterval(() => {
      if (this.id) {
        this.getMetadata(this.id);
      }
    }, 10000);
  }

  ngOnDestroy(): void {
    clearInterval(this.intervalHandler);
  }

  getConfiguration(id: string): void {
    this.instanceService.getDeploymentConfiguration(id).subscribe((data: string) => {
      this.model = data;
      try {
        this.modelYAML = yaml.load(data);
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
    this.instanceService.getMetadata(id).subscribe((data: Metadata) => {
      this.metadata = data;
      this.loadedMetadata = true;
    });
  }

  loadData(): void {
    if (typeof this.id !== 'undefined') {
      this.getConfiguration(this.id);
      this.getInstance(this.id);
      this.getMetadata(this.id);
    }
  }

  isDarkTheme() {
    return this.themeService.isDarkTheme();
  }
}
