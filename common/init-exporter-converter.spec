################################################################################

# rpmbuilder:relative-pack true

################################################################################

%define _posixroot        /
%define _root             /root
%define _bin              /bin
%define _sbin             /sbin
%define _srv              /srv
%define _home             /home
%define _lib32            %{_posixroot}lib
%define _lib64            %{_posixroot}lib64
%define _libdir32         %{_prefix}%{_lib32}
%define _libdir64         %{_prefix}%{_lib64}
%define _logdir           %{_localstatedir}/log
%define _rundir           %{_localstatedir}/run
%define _lockdir          %{_localstatedir}/lock/subsys
%define _cachedir         %{_localstatedir}/cache
%define _spooldir         %{_localstatedir}/spool
%define _crondir          %{_sysconfdir}/cron.d
%define _loc_prefix       %{_prefix}/local
%define _loc_exec_prefix  %{_loc_prefix}
%define _loc_bindir       %{_loc_exec_prefix}/bin
%define _loc_libdir       %{_loc_exec_prefix}/%{_lib}
%define _loc_libdir32     %{_loc_exec_prefix}/%{_lib32}
%define _loc_libdir64     %{_loc_exec_prefix}/%{_lib64}
%define _loc_libexecdir   %{_loc_exec_prefix}/libexec
%define _loc_sbindir      %{_loc_exec_prefix}/sbin
%define _loc_bindir       %{_loc_exec_prefix}/bin
%define _loc_datarootdir  %{_loc_prefix}/share
%define _loc_includedir   %{_loc_prefix}/include
%define _rpmstatedir      %{_sharedstatedir}/rpm-state
%define _pkgconfigdir     %{_libdir}/pkgconfig

################################################################################

%define  debug_package %{nil}

################################################################################

Summary:         Utility for converting init-exporter procfiles from v1 to v2 format
Name:            init-exporter-converter
Version:         0.11.0
Release:         0%{?dist}
Group:           Development/Tools
License:         MIT
URL:             https://github.com/funbox/init-exporter-converter

Source0:         %{name}-%{version}.tar.gz

BuildRoot:       %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:   golang >= 1.14

Provides:        %{name} = %{version}-%{release}

################################################################################

%description
Utility for exporting services described by Procfile to init system.

################################################################################

%prep
%setup -q

%build
export GOPATH=$(pwd)

pushd src/github.com/funbox/%{name}
  %{__make} %{?_smp_mflags} all
popd

%install
rm -rf %{buildroot}

install -dm 755 %{buildroot}%{_bindir}

install -pm 755 src/github.com/funbox/%{name}/%{name} \
                %{buildroot}%{_bindir}/

%clean
rm -rf %{buildroot}

################################################################################

%files
%defattr(-,root,root,-)
%{_bindir}/init-exporter-converter

################################################################################

%changelog
* Thu Jul 30 2020 Anton Novojilov <andy@essentialkaos.com> - 0.11.0-0
- Package ek updated to v12
- Package go-simpleyaml updated to v2

* Mon Aug 05 2019 Anton Novojilov <andyone@fun-box.ru> - 0.10.1-0
- Updated for compatibility with the latest version of ek package

* Thu Jan 10 2019 Anton Novojilov <andyone@fun-box.ru> - 0.10.0-0
- Package ek updated to v10

* Thu Nov 01 2018 Anton Novojilov <andyone@fun-box.ru> - 0.9.1-0
- Rebuilt with the latest procfile parser

* Thu Apr 05 2018 Anton Novojilov <andyone@fun-box.ru> - 0.9.0-1
- Rebuilt with Go 1.10

* Tue Sep 19 2017 Anton Novojilov <andyone@fun-box.ru> - 0.9.0-0
- Improved environment variables parsing in v1

* Fri May 19 2017 Anton Novojilov <andyone@fun-box.ru> - 0.8.0-0
- Migrated to ek.v9
- Added support of multiple file converting
