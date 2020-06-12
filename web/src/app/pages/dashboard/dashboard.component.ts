import { Component, OnDestroy } from '@angular/core';
import { InstanceService } from '../../instances/instance.service';
import 'style-loader!angular2-toaster/toaster.css';
import {
  NbComponentStatus,
  NbGlobalPhysicalPosition,
  NbToastrService
} from '@nebular/theme';

@Component({
  selector: 'ngx-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})

export class DashboardComponent implements OnDestroy {
  public data: any[];
  private intervalHandler;

  constructor(
    private instanceService: InstanceService, 
    private toastrService: NbToastrService) {
    this.getInstances();
    this.intervalHandler = setInterval(() => { this.getInstances() }, 3000);
  }

  getInstances() {
    this.instanceService.getInstances().subscribe((data: any[]) => {
      this.data = data;
    });
  }

  ngOnDestroy() {
    clearInterval(this.intervalHandler);
  }

  delete(id: string) {
    this.instanceService.deleteInstance(id).subscribe(() => {
      this.showToast('success', 'Operation succeeded', 'Deleted Dapr instance with ID ' + id)
    }, error => {
      this.showToast('danger', 'Operation failed', 'Failed to remove Dapr instance with ID ' + id)
    });
  }

  private showToast(type: NbComponentStatus, title: string, body: string) {
    const config = {
      status: type,
      destroyByClick: true,
      duration: 4000,
      hasIcon: true,
      position: NbGlobalPhysicalPosition.TOP_RIGHT,
      preventDuplicates: false,
    };
    const titleContent = title ? `. ${title}` : '';

    this.toastrService.show(
      body,
      titleContent,
      config);
  }
}
