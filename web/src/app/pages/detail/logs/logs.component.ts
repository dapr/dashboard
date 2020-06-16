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
  infoChecked = false;
  debugChecked = false;
  warningChecked = false;
  errorChecked = false;
  fatalChecked = false;

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

  isActive(level: string): boolean {
    if (level === 'info') return this.infoChecked;
    if (level === 'debug') return this.debugChecked;
    if (level === 'warning') return this.warningChecked;
    if (level === 'error') return this.errorChecked;
    if (level === 'fatal') return this.fatalChecked;
    return false;
  }

  toggleInfo(checked: boolean) {
    this.infoChecked = checked;
  }

  toggleDebug(checked: boolean) {
    this.debugChecked = checked;
  }

  toggleWarning(checked: boolean) {
    this.warningChecked = checked;
  }

  toggleError(checked: boolean) {
    this.errorChecked = checked;
  }

  toggleFatal(checked: any) {
    this.fatalChecked = checked;
  }

  goBack(): void {
    this.location.back();
  }
}
