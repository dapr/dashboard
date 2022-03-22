import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { RouterModule } from '@angular/router';
import { DaprComponentsComponent } from './dapr-components.component';

@NgModule({
  imports: [
    CommonModule,
    RouterModule,
    MatTableModule,
    MatIconModule
  ],
  declarations: [
    DaprComponentsComponent
  ],
})
export class DaprComponentsModule { }
