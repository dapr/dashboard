import { TestBed } from '@angular/core/testing';

import { ActorsService } from './actors.service';

describe('ActorsService', () => {
  let service: ActorsService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ActorsService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
