// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {EigenDACertVerifierRouter} from "src/integrations/cert/router/EigenDACertVerifierRouter.sol";

/// @notice For use by rollups to atomically deploy + initialize an immutable CertVerifierRouter (deployed without a proxy).
/// When deployed without a proxy, using this contract is necessary to prevent malicious parties from frontrunning the initialize() transaction and initializing the proxy themselves with byzantine arguments.
contract CertVerifierRouterFactory {
    function deploy(address initialOwner, uint32[] memory initABNs, address[] memory initialCertVerifiers)
        external
        returns (EigenDACertVerifierRouter)
    {
        EigenDACertVerifierRouter router = new EigenDACertVerifierRouter();
        router.initialize(initialOwner, initABNs, initialCertVerifiers);
        return router;
    }
}
