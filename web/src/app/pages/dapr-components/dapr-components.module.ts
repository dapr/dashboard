import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { DaprComponentsComponent } from './dapr-components.component';
import { MatExpansionModule } from '@angular/material/expansion';
import { CommonModule } from '@angular/common';
import { MatTableModule } from '@angular/material/table';

@NgModule({
  imports: [
    CommonModule,
    RouterModule,
    MatExpansionModule,
    MatTableModule,
  ],
  declarations: [
    DaprComponentsComponent,
  ],
})
export class DaprComponentsModule { }
