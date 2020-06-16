import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { DetailComponent } from '../detail/detail.component';
import { MonacoEditorModule } from 'ng-monaco-editor';
import { LogsComponent } from './logs/logs.component';

@NgModule({
  imports: [
    MonacoEditorModule,
    FormsModule,
  ],
  declarations: [
    DetailComponent,
    LogsComponent
  ],
})
export class DetailModule { }
