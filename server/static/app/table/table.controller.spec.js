'use strict';

describe('Controller: TableCtrl', function () {

  // load the controller's module
  beforeEach(module('aptWebApp'));

  var TableCtrl, scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    TableCtrl = $controller('TableCtrl', {
      $scope: scope
    });
  }));

  it('should ...', function () {
    expect(1).toEqual(1);
  });
});
