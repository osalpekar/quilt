const quilt = require('@quilt/quilt');
let infrastructure = require('../../config/infrastructure.js');

let deployment = quilt.createDeployment();
deployment.deploy(infrastructure);

for (let i = 0; i < infrastructure.nWorker; i++) {
    deployment.deploy(new quilt.Service('foo',
        new quilt.Container(
            new quilt.Image('test-custom-image' + i,
                'FROM alpine\n' +
                'RUN echo ' + i + ' > /dockerfile-id\n' +
                'RUN echo $(cat /dev/urandom | tr -dc \'a-zA-Z0-9\' | ' +
                'fold -w 32 | head -n 1) > /image-id'
            ),
            ['tail', '-f', '/dev/null']
        ).replicate(2)
    ));
}
