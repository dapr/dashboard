import { Component, OnInit } from '@angular/core';
import { InstanceService } from '../../instances/instance.service';
import { ActivatedRoute } from '@angular/router';
import * as yaml from 'js-yaml';
import { GlobalsService } from 'src/app/globals/globals.service';
import { ActorsService } from 'src/app/actors/actors.service';
import { Actor } from 'src/app/actors/Actor';

@Component({
  selector: 'ngx-detail',
  templateUrl: './detail.component.html',
  styleUrls: ['./detail.component.scss'],
})
export class DetailComponent implements OnInit {

  private id: string;
  public model: string;
  public modelYAML: any;
  public annotations: string[];
  public options: Object;
  public instance: any;
  public loadedConfiguration: boolean;
  public loadedInstance: boolean;
  public loadedActors: boolean;
  public actors: Actor[];
  public displayedColumns: string[] = ["type", "count"];

  constructor(
    private route: ActivatedRoute,
    private instances: InstanceService,
    private actorsService: ActorsService,
    public globals: GlobalsService,
  ) { }

  ngOnInit() {
    this.loadedConfiguration = false;
    this.loadedInstance = false;
    this.loadedActors = false;
    this.id = this.route.snapshot.params.id;
    this.getConfiguration(this.id);
    this.getInstance(this.id);
    this.getActors(this.id);
    this.options = {
      folding: true,
      minimap: { enabled: true },
      readOnly: false,
      language: 'yaml',
    };
  }

  getConfiguration(id: string): void {
    this.instances.getConfiguration(id).subscribe((data: string) => {
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

  getInstance(id: string) {
    this.instances.getInstance(id).subscribe((data: any) => {
      this.instance = data;
      this.loadedInstance = true;
    });
  }

  getActors(id: string) {
    this.actorsService.getActors(id).subscribe(data => {
      this.actors = data;
      this.loadedActors = true;
    });
  }
}
