import { Component, OnInit } from '@angular/core';
import { InstanceService } from '../../instances/instance.service';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';

@Component({
  selector: 'detail',
  templateUrl: './detail.component.html',
  styleUrls: ['./detail.component.scss']
})
export class DetailComponent implements OnInit {
  private id: string;

  yaml: string[];

  constructor(
    private route: ActivatedRoute, 
    private instances: InstanceService
    ) {  }

  ngOnInit() {
    this.id = this.route.snapshot.params.id;
    this.getYAML(this.id)
  }


  getYAML(id: string): void {
    this.yaml = this.instances.getYAMLArray(id)
  }
}
