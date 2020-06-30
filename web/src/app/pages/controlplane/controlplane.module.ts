import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { ControlPlaneComponent } from './controlplane.component';
import { MatTableModule } from '@angular/material/table';

@NgModule({
  imports: [
    RouterModule,
    MatTableModule,
  ],
  declarations: [
    ControlPlaneComponent,
  ],
})
export class ControlPlaneModule { }
