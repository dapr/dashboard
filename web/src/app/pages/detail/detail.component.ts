import { Component, OnInit } from '@angular/core';
import { InstanceService } from '../../instances/instance.service';
import { ActivatedRoute } from '@angular/router';
import { LogsComponent } from './logs/logs.component';

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
    this.getYAML(this.id);
    this.options = {
      folding: true,
      minimap: { enabled: true },
      readOnly: false,
      language: 'yaml',
    };
  }

  getYAML(id: string): void {
    this.instances.getYAML(id).subscribe((data: string) => {
      this.model = data;
    });
  }
}
