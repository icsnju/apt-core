'use strict';

describe('Controller: SubmitCtrl', function () {

  // load the controller's module
  beforeEach(module('aptWebApp'));

  var SubmitCtrl, scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    SubmitCtrl = $controller('SubmitCtrl', {
      $scope: scope
    });
  }));

  it('should ...', function () {
    expect(1).toEqual(1);
  });
});
