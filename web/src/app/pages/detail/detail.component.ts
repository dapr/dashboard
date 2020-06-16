import { Component, OnInit } from '@angular/core';
import { InstanceService } from '../../instances/instance.service';
import { ActivatedRoute } from '@angular/router';
<<<<<<< HEAD
import { LogsComponent } from './logs/logs.component';
=======
>>>>>>> develop

@Component({
  selector: 'ngx-detail',
  templateUrl: './detail.component.html',
  styleUrls: ['./detail.component.scss'],
})
export class DetailComponent implements OnInit {
  private id: string;

  model: string;
  options: Object;

  constructor(
    private route: ActivatedRoute,
    private instances: InstanceService) {}

  ngOnInit() {
    this.id = this.route.snapshot.params.id;
<<<<<<< HEAD
    this.getYAML(this.id);
=======
    this.getConfiguration(this.id);
>>>>>>> develop
    this.options = {
      folding: true,
      minimap: { enabled: true },
      readOnly: false,
      language: 'yaml',
    };
  }

<<<<<<< HEAD
  getYAML(id: string): void {
    this.instances.getYAML(id).subscribe((data: string) => {
=======
  getConfiguration(id: string): void {
    this.instances.getConfiguration(id).subscribe((data: string) => {
>>>>>>> develop
      this.model = data;
    });
  }
}
