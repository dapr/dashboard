import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { InstanceService } from '../../instances/instance.service';

import 'style-loader!angular2-toaster/toaster.css';
import {
  NbComponentStatus,
  NbGlobalPhysicalPosition,
  NbToastrService,
} from '@nebular/theme';
@Component({
  selector: 'logs',
  templateUrl: './logs.component.html'
})
export class LogsComponent implements OnInit {
  logs: string;
  id: string;
  constructor(private route: ActivatedRoute, private instances: InstanceService, private toastrService: NbToastrService) { }

  ngOnInit() {
    this.id = this.route.snapshot.params.id;
    this.getLogs(false);
  }

  getLogs(showMessage: boolean) {
    this.instances.getLogs(this.id).subscribe((logData: string) => {
      this.logs = logData;
    });

    if (showMessage) {
      this.showToast("info")
    }
  }

  private showToast(type: NbComponentStatus) {
    const config = {
      status: type,
      destroyByClick: true,
      duration: 4000,
      hasIcon: true,
      position: NbGlobalPhysicalPosition.TOP_RIGHT,
      preventDuplicates: false,
    };
    this.toastrService.show(
      "Logs refreshed",
      "Status",
      config);
  }
}
