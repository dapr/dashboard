import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';
import { InstanceService } from '../../../instances/instance.service';
import { Log } from './log';

@Component({
  selector: 'ngx-logs',
  templateUrl: './logs.component.html',
})

export class LogsComponent implements OnInit {
  logs: Log[];
  id: string;

  constructor(
    private route: ActivatedRoute,
    private instances: InstanceService,
    private location: Location) { }

  ngOnInit() {
    this.id = this.route.snapshot.params.id;
    this.getLogs(false);
  }

  getLogs(showMessage: boolean) {
    this.logs = this.instances.getLogsArray(this.id);

    if (showMessage) {
    }
  }

  goBack(): void {
    this.location.back();
  }
}
