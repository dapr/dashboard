import { Component, OnInit } from '@angular/core';
import { InstanceService } from 'src/app/instances/instance.service';
import { ActivatedRoute } from '@angular/router';
import * as yaml from 'js-yaml';
import { GlobalsService } from 'src/app/globals/globals.service';
import { Metadata, Instance } from 'src/app/types/types';

@Component({
  selector: 'app-detail',
  templateUrl: 'detail.component.html',
  styleUrls: ['detail.component.scss'],
})
export class DetailComponent implements OnInit {

  private id: string;
  public model: string;
  public modelYAML: any;
  public annotations: string[];
  public options: object;
  public instance: Instance;
  public loadedConfiguration: boolean;
  public loadedInstance: boolean;
  public loadedMetadata: boolean;
  public metadata: Metadata[];
  public metadataDisplayedColumns: string[] = ['type', 'count'];

  constructor(
    private route: ActivatedRoute,
    private instanceService: InstanceService,
    public globals: GlobalsService,
  ) { }

  ngOnInit(): void {
    this.loadedConfiguration = false;
    this.loadedInstance = false;
    this.loadedMetadata = false;
    this.id = this.route.snapshot.params.id;
    this.getConfiguration(this.id);
    this.getInstance(this.id);
    this.getMetadata(this.id);
    this.options = {
      folding: true,
      minimap: { enabled: true },
      readOnly: false,
      language: 'yaml',
    };
  }

  getConfiguration(id: string): void {
    this.instanceService.getConfiguration(id).subscribe((data: string) => {
      this.model = data;
      try {
        this.modelYAML = yaml.safeLoad(data);
        this.annotations = Object.keys(this.modelYAML.metadata.annotations);
        this.loadedConfiguration = true;
      } catch (e) {
        this.modelYAML = {};
      }
    });
  }

  getInstance(id: string): void {
    this.instanceService.getInstance(id).subscribe((data: Instance) => {
      this.instance = data;
      this.loadedInstance = true;
    });
  }

  getMetadata(id: string): void {
    this.instanceService.getMetadata(id).subscribe((data: Metadata[]) => {
      this.metadata = data;
      this.loadedMetadata = true;
    });
  }
}
