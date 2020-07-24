import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'filter'
})
export class FilterPipe implements PipeTransform {
  transform(items: any[], searchKey: string, fieldName: string): any[] {
    if (!items) { return []; }
    if (!searchKey) { return items; }
    searchKey = searchKey.toLowerCase();
    return items.filter(item => {
      if (item && item[fieldName]) {
        return item[fieldName].toLowerCase().includes(searchKey);
      }
      return false;
    });
   }
}
