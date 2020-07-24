import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { DaprComponentsComponent } from './dapr-components.component';
import { CommonModule } from '@angular/common';
import { MatTableModule } from '@angular/material/table';

@NgModule({
  imports: [
    CommonModule,
    RouterModule,
    MatTableModule,
  ],
  declarations: [
    DaprComponentsComponent
  ],
})
export class DaprComponentsModule { }
