import { Pipe, PipeTransform } from '@angular/core';
import { Log } from 'src/app/types/types';

@Pipe({
  name: 'filter'
})
export class FilterPipe implements PipeTransform {
  transform(items: Log[], searchKey: string, showFiltered: boolean): Log[] {
    if (!items) { return []; }
    if (!searchKey || !showFiltered) { return items; }
    return items.filter(item => {
      if (item && item.content) {
        return item.content.includes(searchKey);
      }
      return false;
    });
  }
}
