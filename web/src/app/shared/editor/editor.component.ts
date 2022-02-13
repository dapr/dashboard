import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { ThemeService } from '../../theme/theme.service';
import { YamlViewerOptions } from '../../types/types';

@Component({
  selector: 'app-editor',
  templateUrl: 'editor.component.html'
})
export class EditorComponent implements OnInit {

  @Input() model = '';
  @Input() language = '';
  @Output() modelChange = new EventEmitter<string>();

  options: YamlViewerOptions = {
    folding: true,
    minimap: { enabled: true },
    readOnly: true,
    language: this.language || 'yaml',
    contextmenu: false,
    scrollBeyondLastLine: false,
    lineNumbers: false as any,
    theme: this.isDarkTheme() ? 'vs-dark' : 'vs'
  };

  constructor(
    private themeService: ThemeService,
  ) { }

  ngOnInit() {
    this.model = this.model.repeat(10);
    this.themeService.themeChanged.subscribe((newTheme: string) => {
      this.options = {
        ...this.options,
        theme: newTheme.includes('dark') ? 'vs-dark' : 'vs'
      };
    });
  }

  isDarkTheme(): boolean {
    return this.themeService.getTheme().includes('dark');
  }
}
