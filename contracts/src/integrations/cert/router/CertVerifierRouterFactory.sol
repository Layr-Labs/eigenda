// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {EigenDACertVerifierRouter} from "src/integrations/cert/router/EigenDACertVerifierRouter.sol";

/// @notice For use by rollups to deploy a CertVerifierRouter atomically without a proxy.
contract CertVerifierRouterFactory {
    function deploy(address initialOwner, uint32[] memory initialRBNs, address[] memory initialCertVerifiers)
        external
        returns (EigenDACertVerifierRouter)
    {
        EigenDACertVerifierRouter router = new EigenDACertVerifierRouter();
        router.initialize(initialOwner, initialRBNs, initialCertVerifiers);
        return router;
    }
}
