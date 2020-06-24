import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { InstanceService } from '../../../instances/instance.service';
import { Log } from './log';

@Component({
  selector: 'ngx-logs',
  templateUrl: './logs.component.html',
  styleUrls: ['./logs.component.scss'],
})

export class LogsComponent implements OnInit {
  logs: Log[];
  id: string;
  info: boolean;
  debug: boolean;
  warning: boolean;
  error: boolean;
  fatal: boolean;

  constructor(
    private route: ActivatedRoute,
    private instances: InstanceService) { }

  ngOnInit() {
    this.id = this.route.snapshot.params.id;
    this.getLogs(false);
  }

  getLogs(showMessage: boolean) {
    this.logs = this.instances.getLogsArray(this.id);

    if (showMessage) {
    }
  }

  isActive(level: string): boolean {
    if (level === "info") return this.info;
    if (level === "debug") return this.debug;
    if (level === "warning") return this.warning;
    if (level === "error") return this.error;
    if (level === "fatal") return this.fatal;
    return false;
  }
}
