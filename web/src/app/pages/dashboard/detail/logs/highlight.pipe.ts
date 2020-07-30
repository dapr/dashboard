import { PipeTransform, Pipe } from '@angular/core';

@Pipe({
    name: 'highlight'
})
export class HighlightPipe implements PipeTransform {
    transform(item: string, searchKey: string): string {
        return item.replace(searchKey, (match) => `<span class="highlighted">${match}</span>`);
    }
}
