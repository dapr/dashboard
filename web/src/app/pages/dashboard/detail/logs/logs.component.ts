import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { InstanceService } from 'src/app/instances/instance.service';
import { LogStreamService } from 'src/app/logstream/logstream.service';
import { Log } from 'src/app/types/types';

@Component({
  selector: 'app-logs',
  templateUrl: 'logs.component.html',
  styleUrls: ['logs.component.scss'],
})

export class LogsComponent implements OnInit, OnDestroy {

  public id: string;
  public containers: string[];
  public showFiltered = true;
  public filterValue = '';
  public containerValue = 'daprd';
  public since: number;
  public sinceUnit = '';
  public logs: Map<string, Log[]> = new Map();

  constructor(
    private route: ActivatedRoute,
    private instances: InstanceService,
    private logStream: LogStreamService,
  ) { }

  ngOnInit(): void {
    this.id = this.route.snapshot.params.id;
    this.containers = ['daprd'];
    this.instances.getContainers(this.id).subscribe(containers => {
      this.containers = containers;
      containers.forEach(container => {
        const containerLogs = new Array<Log>();
        this.logs.set(container, containerLogs);
        this.logStream.startStream(this.id, container).subscribe(logRecord => {
           containerLogs.push(logRecord);
        });
      });
    });
  }

  ngOnDestroy(): void {
    this.containers.forEach(container => {
      this.logStream.endStream(this.id, container);
    });
  }

}
